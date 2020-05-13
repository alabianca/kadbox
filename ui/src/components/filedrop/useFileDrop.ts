import {useEffect, useState} from "react";

export const useFileDrop = () => {
    const [dragActive, setDragActive] = useState(false);
    const [dragCounter, setDragCounter] = useState(0);
    const [fileList, setFileList] = useState([]);


    useEffect(() => {
        if (dragCounter > 0) {
            setDragActive(true)
        } else {
            setDragActive(false);
        }
    }, [dragCounter]);

    const onDragEnter = (e: any) => {
        e.stopPropagation();
        e.preventDefault();
        setDragCounter(dragCounter + 1);

    };

    const onDragLeave = (e: any) => {
        setDragCounter(dragCounter - 1);

    };

    const onDragOver = (e: any) => {
        e.stopPropagation();
        e.preventDefault();

    };


    const onDrop = (e: any) => {
        e.stopPropagation();
        e.preventDefault();
        setDragCounter(0);
        setFileList(e.dataTransfer.files);

    };

    return {
        dragActive,
        fileList,
        onDragEnter,
        onDragLeave,
        onDragOver,
        onDrop,
    }
};
