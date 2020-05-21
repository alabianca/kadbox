import React, {useEffect} from 'react';
import {Switch, Route} from 'react-router-dom';
import Lib from '../../lib';
import { IpcEventBus } from "../../lib/eventBus";

import './content.css'
import FileDrop from "../filedrop/filedrop";
import FileList from "../filelist/filelist";
import Downloads from "../downloads/downloads";

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
                    <div className='content-container'>
                        <FileList/>
                    </div>
                </Route>
                <Route path="/downloads">
                    <div className='content-container'>
                        <Downloads/>
                    </div>
                </Route>
            </Switch>
        </div>
    )
};

export default Content;