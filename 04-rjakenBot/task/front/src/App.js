import './App.css';
import {useState,} from "react";
import Dropdown from 'react-dropdown';

function App() {
    const [input, setInput] = useState()
    const [isValid, setIsValid] = useState(true)
    const [validUrl, setValidUrl] = useState(true)
    const [img, setImg] = useState();
    const [method, setMethod] = useState("GET")
    const [urlError, setUrlError] = useState()

    function handleChange(event) {
        setValidUrl(true)
        const regex = new RegExp(/(?:http|https):\/\/(?:[\w\-]+\.)*[\w\-]+(?:\:\d{1,5})?(?:\/[\w\-]*)*(?:\?.*)?/gi);
        if (event.match(regex) || event === "") {
            setIsValid(true)
            setInput(event)
        } else {
            setIsValid(false)
        }
    }

    function dropdownChange(option) {
        setMethod(option["value"])
    }

    async function fetchBackend(e) {
        const requestOptions = {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({"pictureLink": `${input}`, "method": `${method}`})
        };
        const res = await fetch('/api/generate-meme', requestOptions)
        if (!res.ok) {
            const error = await res.json();
            setUrlError(error["error"])
            setValidUrl(false)
            setImg()
        } else {
            setValidUrl(true)
            const imageBlob = await res.blob();
            const imageObjectURL = URL.createObjectURL(imageBlob);
            setImg(imageObjectURL)
        }

        // https://thumb.mp-farm.com/62531672/preview.jpg
        // console.log(img)
    }


    return (
        <div className="grandparent">
            <div className="parent">
                <header className="service-name">МЕГАДЭВАЙС 3000 v.2</header>
                <div className="fields">

                    <Dropdown options={["GET", "POST"]}
                              onChange={dropdownChange}
                              placeholder="Method"
                              className="dropdown"
                              controlClassName='dropdownControl'
                              menuClassName='dropdownMenu'
                              placeholderClassName='dropdownPlaceholder'/>

                    <div className="div-input">
                        <input placeholder="Введи URL картинки, чтобы получить classique meme)" type="text" className="input"
                               onChange={e => handleChange(e.target.value)}/>
                        {!isValid && (
                            <span className="spanchik">Введите валидный урл :( </span>
                        )}
                        {!validUrl && (
                            <span className="spanchik">{urlError}</span>
                        )}
                        <div className="imaga">
                            {/*{!validUrl && (*/}
                            {/*    <div className="url"> {urlError}</div>*/}
                            {/*)}*/}
                            <img src={img} alt="" className="image"/>
                        </div>
                    </div>
                </div>
                <button onClick={fetchBackend} className="button" disabled={!isValid}>Палучить мем</button>
            </div>
        </div>
    );
}

export default App;
