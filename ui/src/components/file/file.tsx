import React from 'react';
import './file.css'

interface FileProps {
    name: string,
}

const File = (props: FileProps) => {
    const displayHash = `${props.name.substring(0, 8)}...${props.name.substring(props.name.length - 8)}`
    return (
        <div className='file'>
            <span className='scheme'>kadbox://</span><span className='hash'>{displayHash}</span>
        </div>
    )
};

export default File;