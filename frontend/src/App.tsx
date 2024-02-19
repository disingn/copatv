import React, { useState } from 'react'
import './App.css'
import { GetVersion,OpenUrl } from '../wailsjs/go/main/App';

function App() {
    const [resultText, setResultText] = useState<string>("未授权🥰");
    const [name, setName] = useState<string>("");
    const [auth, setAuth] = useState<string>('');
    const updateAuth = (event: React.ChangeEvent<HTMLInputElement>) => {
        setAuth(event.target.value);
    }
    const updateName = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName(event.target.value);
    }
    const [tool, setTool] = useState<string>('jetbrains');
    const updateTool = (event: React.ChangeEvent<HTMLSelectElement>) => {
        setTool(event.target.value);
    }
    const updateResultText = (result: string) => setResultText(result);
    // function greet() {
    //     Greet(auth).then(updateResultText);
    // }
    function getVersion() {
        if(auth === '') {
            setResultText("请输入授权码")
            return
        }
        if(tool === '') {
            setResultText("请选择工具")
            return
        }
        if (name === '') {
            setResultText("请输入用户名")
            return
        }
        GetVersion(auth, tool,name).then(updateResultText);
        setResultText("")
    }
    function OpenBrowser() {
        OpenUrl("https://www.aiu.im").then()
    }
    return (
        <div className="flex justify-center items-center min-h-screen bg-white">
            <div className="space-y-4 w-full max-w-xs mx-auto">
                <div className="flex flex-col w-full">
                    <span className="mb-2">授权码:</span>
                    <input
                        type="text"
                        placeholder="输入你的授权码"
                        className="input input-bordered w-full"
                        value={auth}
                        onChange={updateAuth}
                    />
                </div>
                <div className="flex flex-col w-full">
                    <span className="mb-2">bio:</span>
                    <input
                        type="text"
                        placeholder="输入你的用户名"
                        className="input input-bordered w-full"
                        value={name}
                        onChange={updateName}
                    />
                </div>
                <div className="flex flex-col w-full">
                    <span className="mb-2">工具:</span>
                    <select
                        className="select select-bordered w-full"
                        value={tool}
                        onChange={updateTool}
                    >
                        <option value="" disabled>选择一个</option>
                        <option value="jetbrains">jetbrains</option>
                        <option value="vscode">vscode</option>
                    </select>
                </div>

                {/*<label className="form-control w-full">*/}
                {/*    <div className="label">*/}
                {/*        <span className="label-text">Your bio</span>*/}
                {/*    </div>*/}
                {/*    <textarea*/}
                {/*        className="textarea textarea-bordered h-24 w-full"*/}
                {/*        value={resultText}*/}
                {/*        onChange={(e) => setResultText(e.target.value)}*/}
                {/*    />*/}
                {/*</label>*/}
                <button className="btn btn-primary w-full" onClick={getVersion}>授权</button>
                <p className="text-gray-500 text-xs text-center">{resultText}</p>
                <a
                    onClick={OpenBrowser}
                    className="block w-full mt-6 text-center text-gray-500  hover:text-blue-700 link"
                >
                    前往官网了解更多信息～
                </a>

            </div>

        </div>
    );
}


export default App
