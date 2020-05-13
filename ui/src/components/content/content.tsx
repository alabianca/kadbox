import React from 'react';

import {Switch, Route} from 'react-router-dom'

import './content.css'
import FileDrop from "../filedrop/filedrop";

const Content = () => {
    return (
        <div className='content'>
            <Switch>
                <Route path="/share">
                    <div className='filedrop-container'>
                        <FileDrop/>
                    </div>

                </Route>
            </Switch>
        </div>
    )
};

export default Content;