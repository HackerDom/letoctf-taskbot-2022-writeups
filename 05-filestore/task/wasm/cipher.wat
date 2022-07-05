(module
    ;; data layout: encryption key | PI generated round keys | PI generated sboxes
    (import "imports" "mem" (memory 1))

    (func $log32 (import "imports" "log") (param i32))
    (func $log64 (import "imports" "log") (param i64))

    (global $keySize i32 (i32.const 16))
    (global $sizeofI32 i32 (i32.const 4))

    (func (export "encryptBlock") (param $blockOff i32)
        (local $left i32)
        (local $right i32)

        local.get $blockOff
        i32.load

        local.get $blockOff
        global.get $sizeofI32
        i32.add
        i32.load

        call $encryptBlockInternal
        local.set $right
        local.set $left

        local.get $blockOff
        local.get $left
        i32.store

        local.get $blockOff
        global.get $sizeofI32
        i32.add
        local.get $right
        i32.store
    )

    (func $encryptBlockInternal 
        (param $left i32) 
        (param $right i32) 
        (result i32) (result i32)

        (local $round i32)
        (local $p16offset i32)

        i32.const 0
        local.set $round

        (loop
            ;; load round key
            local.get $round
            i32.const 4
            i32.mul
            global.get $keySize
            i32.add
            i32.load

            ;; left = left ^ pn
            local.get $left
            i32.xor
            local.set $left

            ;; new_left = ((S0[a] + S1[b] ^ S2[c]) + S3[d] & 0xffffffff) ^ R
            local.get $left
            call $f
            local.get $right
            i32.xor

            ;; right = left
            local.get $left
            local.set $right

            ;; left = new_left
            local.set $left

            ;; round += 1
            local.get $round
            i32.const 1
            i32.add
            local.set $round

            ;; if (round < 16) goto loop;
            local.get $round
            i32.const 16 ;; round count
            i32.lt_u
            br_if 0
        )
        
        i32.const 16
        global.get $sizeofI32
        i32.mul
        global.get $keySize
        i32.add
        local.set $p16offset

        ;; right ^ p17
        local.get $right
        local.get $p16offset
        global.get $sizeofI32
        i32.add
        i32.load
        i32.xor

        ;; left ^ p16
        local.get $left
        local.get $p16offset
        i32.load
        i32.xor
    )

    (func $getByte (param $num i32) (param $idx i32) (result i32)
        (local $shiftBy i32)

        global.get $sizeofI32
        i32.const 1
        i32.sub ;; index of last element 
        local.get $idx
        i32.sub ;; little endian
        i32.const 8 ;; num of bits
        i32.mul
        local.set $shiftBy

        local.get $num
        i32.const 255 ;; 0xff
        local.get $shiftBy
        i32.shl
        i32.and
        local.get $shiftBy
        i32.shr_u
    )

    (func $sbox (param $src i32) (param $idx i32) (result i32)
        (local $sboxesOff i32)

        i32.const 18
        global.get $sizeofI32
        i32.mul
        global.get $keySize
        i32.add
        local.set $sboxesOff

        local.get $src
        local.get $idx
        call $getByte ;; number of sbox value
        global.get $sizeofI32
        i32.mul ;; byte index of sbox value

        ;; count offset
        i32.const 1024 ;; size of sbox (256 * 4)
        local.get $idx
        i32.mul
        local.get $sboxesOff
        i32.add

        ;; offset + gotten byte index of sbox value
        i32.add
        ;; load sbox mapping
        i32.load
    )

    (func $f (param $src i32) (result i32)
        ;; getting S0[a]
        local.get $src
        i32.const 0
        call $sbox
        i64.extend_i32_s

        ;; getting S1[b]
        local.get $src
        i32.const 1
        call $sbox
        i64.extend_i32_s

        ;; S0[a] + S1[b] 
        i64.add

        ;; getting S2[c]
        local.get $src
        i32.const 2
        call $sbox
        i64.extend_i32_s

        ;; S0[a] + S1[b] ^ S2[c]
        i64.xor

        ;; getting S3[d]
        local.get $src
        i32.const 3
        call $sbox
        i64.extend_i32_s

        ;; (S0[a] + S1[b] ^ S2[c]) + S3[d]
        i64.add

        ;; (S0[a] + S1[b] ^ S2[c]) + S3[d] & 0xffffffff
        i64.const 4294967295
        i64.and

        i32.wrap_i64
    )

    (func (export "init")
        call $xorKeys
        call $encryptKeys
    )

    (func $encryptKeys 
        (local $i i32)
        (local $left i32)
        (local $right i32)

        global.get $keySize
        local.set $i

        (loop 
            local.get $left
            local.get $right
            call $encryptBlockInternal
            local.set $right
            local.set $left

            local.get $i
            local.get $left
            i32.store

            local.get $i
            global.get $sizeofI32
            i32.add
            local.get $right
            i32.store

            ;; i += 8
            local.get $i
            global.get $sizeofI32
            i32.const 2
            i32.mul
            i32.add
            local.set $i

            ;; if (i < 4184) goto loop
            local.get $i
            global.get $keySize
            i32.const 4168 ;; 18 round keys, 256 * 4 sboxes
            i32.add
            i32.lt_u

            br_if 0
        )
    )
    
    (func $xorKeys 
        (local $i i32)
        (local $currKey i32)
        (local $j i32)
        
        global.get $keySize
        local.set $i

        (loop
            local.get $i
            i32.load
            local.set $currKey

            i32.const 0
            local.set $j

            (loop
                ;; position of storing result byte
                ;; i32 little endian 
                local.get $i
                local.get $j
                i32.sub
                i32.const 3
                local.get $j
                i32.sub
                i32.add

                ;; get byte of key
                local.get $i
                global.get $keySize
                i32.const 1
                i32.sub
                i32.and
                i32.load8_u

                ;; get byte of round key
                local.get $currKey
                local.get $j
                call $getByte

                i32.xor

                i32.store8

                ;; i++
                local.get $i
                i32.const 1
                i32.add
                local.set $i

                ;; j++
                local.get $j
                i32.const 1
                i32.add
                local.set $j

                ;; if (j < sizeofI32) goto loop
                local.get $j
                global.get $sizeofI32
                i32.lt_u

                br_if 0
            )

            ;; if (i < 88) goto loop
            local.get $i
            global.get $keySize
            i32.const 72 ;; 18 round keys
            i32.add
            i32.lt_u

            br_if 0
        )
    )
)