import React, {useEffect, useState} from "react";

interface TagProps {
    text: string,
    type: string,
    isLight?: boolean,
}

interface TimeoutTagProps {
    visible: boolean,
    timeout: number,
    text: string,
    type: string,
    isLight?: boolean,
    onFadeOut: () => any
}

const Tag = (props: TagProps) => {
    return (
        <>
            <span className={`tag ${props.type} ${props.isLight ? 'is-light': ''}`}>{props.text}</span>
        </>
    )
}

export const TagWithTimeout = (props: TimeoutTagProps) => {
    const [isVisible, setIsVisible] = useState(props.visible)

    useEffect(() => {
        setIsVisible(props.visible);
        const id = setTimeout(() => {
            setIsVisible(false);
            props.onFadeOut();
        }, props.timeout || 2000);

        return function cleanup() {
            clearTimeout(id)
        }
    }, [props.visible])

    return (
        <>
            {isVisible ?  <Tag text={props.text} type={props.type} isLight={props.isLight} /> : ''}
        </>
    )
}

export default Tag