import { Alert, Box, Checkbox, Input, Snackbar, Typography } from '@mui/material';
import React, { useContext, useState } from 'react';
import s from './UploadPage.module.css';
import UploadFileIcon from "@mui/icons-material/UploadFile";
import NavBar, { Page } from '../../components/NavBar/NavBar';
import { authContext } from '../../auth/context';
import { Navigate } from 'react-router-dom';
import { parse as parseUuid } from 'uuid';

export function UploadPage() {
    enum uploadCodes {
        empty = '',
        success = 'successful',
    }
    const [code, changeCode] = useState(uploadCodes.empty); 
    const [encrypted, changeEncrypted] = useState(false);
    const ctx = useContext(authContext);

    if (!ctx.loggedIn) {
        return <Navigate to="/" />
    }

    const onChangeFn = async (e: any) => {
        const data = new FormData();
        data.append('encrypted', encrypted.toString());

        if (encrypted) {
            const encryptedFile = await encrypt(e.target.files[0], ctx.userId);
            data.append('file', encryptedFile);
        } else {
            data.append('file', e.target.files[0]);
        }

        const r = await fetch('/api/upload', {
            method: 'PUT',
            body: data,
        });
        if (r.ok) {
            changeCode(uploadCodes.success);
            return;
        }

        const errorReason = (await r.json()).response;
        changeCode(errorReason);
    };

    const onClose = () => {
        changeCode(uploadCodes.empty);
    };

    return (
        <React.Fragment>
            <Box
                display="flex"
                height="100%"
                flexDirection="column"
            >
                <NavBar page={Page.Upload} />
                <Box
                    display="flex"
                    justifyContent="center"
                    alignItems="center"
                    flexDirection="column"
                    flexGrow="1"
                >
                    <label htmlFor="file">
                        <div className={s.fileInputHidden} >
                            <Input id="file"type="file" onChange={onChangeFn} />
                        </div>
                        <div className={s.uploadElement}>
                            <UploadFileIcon sx={{ fontSize: 100 }} />
                            <Typography align="center">Upload file</Typography>
                        </div>
                    </label>
                    <Box display="flex" alignItems="center">
                        <Checkbox value={encrypted} onChange={e => changeEncrypted(e.target.checked)} />
                        <Typography>encrypt</Typography>
                    </Box>
                </Box>
                <Snackbar
                    open={code === uploadCodes.success}
                    autoHideDuration={4000}
                    onClose={onClose}
                >
                    <Alert severity="success">File uploaded</Alert>
                </Snackbar>
                <Snackbar
                    open={code !== uploadCodes.success && code !== uploadCodes.empty}
                    autoHideDuration={4000}
                    onClose={onClose}
                >
                    <Alert severity="error">{code}</Alert>
                </Snackbar>
            </Box>
        </React.Fragment>
    );
}

async function encrypt(file: File, userId: string): Promise<File> {
    const fileBuf = await readFile(file);

    const psSize = 4168;
    const trashSize = 832;

    const key = parseUuid(userId);

    const psbinResp = await fetch('/api/wasm/ps.bin');
    const psbin = new Uint8Array(await psbinResp.arrayBuffer(), 0, psSize);

    const memory = new WebAssembly.Memory({ initial: 1 });
    const data = new Uint8Array(memory.buffer, 0, key.length + psSize + trashSize + fileBuf.length);
    
    for (let i = 0; i < key.length; i++) {
        data[i] = key[i];
    }

    for (let i = 0; i < psSize; i++) {
        data[key.length + i] = psbin[i];
    }

    for (let i = 0; i < trashSize; i++) {
        data[key.length + psSize + i] = 0;
    }

    const blockOff = key.length + psSize + trashSize;

    for (let i = 0; i < fileBuf.length; i++) {
        data[blockOff + i] = fileBuf[i];
    }

    const imports = {
        log: console.log,
        mem: memory,
    };
    const cipher = await WebAssembly.instantiateStreaming(fetch('/api/wasm/cipher.wasm'), { imports });
    (cipher.instance.exports as any).init();

    for (let off = blockOff; off < data.length; off += 8) {
        (cipher.instance.exports as any).encryptBlock(off);
    }

    const sl = data.slice(blockOff, blockOff + fileBuf.length);
    return new File([new Blob([sl])], file.name);
}

function readFile(file: File): Promise<Uint8Array> {
    const reader = new FileReader();
    reader.readAsArrayBuffer(file);
    
    return new Promise((resolve, reject) => {
        reader.onload = () => {
            const result = new Uint8Array(reader.result as ArrayBuffer);
            const len = Math.ceil(result.byteLength / 8) * 8;
            const buf = new Uint8Array(new ArrayBuffer(len));

            for (let i = 0; i < result.byteLength; i++) {
                buf[i] = result[i];
            }
            
            resolve(buf);
        };
        reader.onerror = () => {
            reject(reader.result);
        };
    });
} 
