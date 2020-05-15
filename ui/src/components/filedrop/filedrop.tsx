import React, {useEffect, useState} from 'react';
import './filedrop.css'
import {useFileDrop} from './useFileDrop';
import Lib from '../../lib'

const FileDrop = () => {
    const {dragActive, fileList, onDragEnter, onDragLeave, onDragOver, onDrop} = useFileDrop();
    const fmgr = Lib.getConnectionManager().FileManager();

    useEffect(() => {
        console.log(fileList)
        for (let f of fileList) {
            fmgr.addObject(f.path);
        }

    }, [fileList])

    return (
        <div
            onDrop={onDrop}
            id='drag'
            onDragOver={onDragOver}
            onDragLeave={onDragLeave}
            onDragEnter={onDragEnter}
            className={`filedrop ${dragActive ? 'active': ''}`}>
            <p id="drag-text">Drag Files or Folders to Share them in the Network</p>
        </div>
    )
};

export default FileDrop;