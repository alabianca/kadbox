import { useState, useEffect } from 'react';
import Lib from "../lib";

export const useFiles = (initialFiles: string[] = []) => {
    const [files, setFiles] = useState(initialFiles);
    const fmgr = Lib.getConnectionManager().FileManager();

    useEffect(() => {
        fmgr.getFiles((f) => {
            setFiles(f)
        })
    })

    return {
        files,
    }
}