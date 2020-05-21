import React from 'react';
import './file.css'

interface FileProps {
    name: string,
    onCopy: (url: string) => any,
}

const File = (props: FileProps) => {
    const displayHash = `${props.name.substring(0, 8)}...${props.name.substring(props.name.length - 8)}`
    let copyInput: HTMLInputElement;

    const copyToClipboard = () => {
        copyInput.select();
        document.execCommand('copy');
        props.onCopy(`kadbox://${props.name}`)
    }

    return (
        <div className='file'>
            <span>kadbox://{displayHash}</span>
            <button onClick={() => copyToClipboard()}>Copy</button>
            <input style={{display: "inline-block", marginLeft: "999px"}}  readOnly={true} value={props.name} ref={(input) => copyInput = input}  type='text'/>
        </div>
    )
};

export default File;