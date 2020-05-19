import React, {useEffect} from 'react';
import {Switch, Route} from 'react-router-dom';
import Lib from '../../lib';
import { IpcEventBus } from "../../lib/eventBus";

import './content.css'
import FileDrop from "../filedrop/filedrop";
import FileList from "../filelist/filelist";

const Content = () => {


    return (
        <div className='content'>
            <Switch>
                <Route path="/share">
                    <div className='filedrop-container'>
                        <FileDrop/>
                    </div>
                </Route>
                <Route path="/files">
                    <div className='files-container'>
                        <FileList/>
                    </div>
                </Route>
            </Switch>
        </div>
    )
};

export default Content;