import React from 'react'
import {useFiles} from "../../hooks/useFiles";
import File from "../file/file";

const FileList = () => {
    const { files } = useFiles()
    return (
        <div>
            {
                files.map((file, index) => <File key={index} name={file} />)
            }
        </div>
    )
}

export default FileList;