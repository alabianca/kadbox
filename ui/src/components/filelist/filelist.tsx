import React, {useState} from 'react'
import {useFiles} from "../../hooks/useFiles";
import File from "../file/file";
import {TagWithTimeout} from "../tag/tag";

const FileList = () => {
    const { files } = useFiles();
    const [copyTagVisible, setCopyTagVisible] = useState(0)

    const onUrlCopied = (url: string) => {
        setCopyTagVisible(copyTagVisible + 1)
    }

    console.log('hello?')

    return (
        <div className='file-list'>
            {
                files.map((file, index) => <File onCopy={(url) => onUrlCopied(url)} key={index} name={file} />)
            }
            <TagWithTimeout onFadeOut={() => setCopyTagVisible(0)} visible={copyTagVisible > 0} timeout={2000} text={'Copied!'} type={'is-success'} isLight={true} />
        </div>
    )
}

export default FileList;