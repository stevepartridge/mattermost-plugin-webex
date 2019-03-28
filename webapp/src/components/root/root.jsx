// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import React from 'react';
import PropTypes from 'prop-types';
import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

import {id as pluginId} from '../../manifest';

import mmLogoURL from './../../images/mm-logo.png';
import webexLogoURL from './../../images/webex-logo.png';

const {Modal} = window.ReactBootstrap;

export default class Root extends React.PureComponent {
    static propTypes = {

        theme: PropTypes.object.isRequired,
        siteUrl: PropTypes.string.isRequired,
        visible: PropTypes.bool.isRequired,
        authenticated: PropTypes.object.isRequired,
        actions: PropTypes.shape({
            closeRootModal: PropTypes.func.isRequired,
        }).isRequired,

    }

    handleClose = () => {
        this.props.actions.closeRootModal();
    };

    startOAuthConnectFlow = async () => {
        window.open(`${this.props.siteUrl}/plugins/${pluginId}/oauth2/connect`);
    }

    render() {
        var posWidth = ((window.innerWidth - 200) + 'px');
        var style = getStyle(posWidth, this.props.theme);

        return (
            <div style={style.modelCont}>

                <Modal
                    show={this.props.visible}
                    onHide={this.handleClose}
                    style={style.modal}
                >
                    <Modal.Header
                        closeButton={true}
                        style={style.header}
                    ></Modal.Header>

                    <Modal.Body style={style.body}>
                        <div >
                            <div style={style.logosConnect}>
                                <img
                                    style={style.mmLogo}
                                    src={mmLogoURL}
                                    className='img-responsive img-circle center-block pull-left'
                                    width='100'
                                />
                                <img
                                    style={style.webexLogo}
                                    src={webexLogoURL}
                                    className='img-responsive img-circle center-block pull-right'
                                    width='100'
                                />
                            </div>
                            <div style={style.bodyText}>
                                <span >
                                    <i
                                        style={style.connectArrow}
                                        className='fa fa-arrow-right fa-2x'
                                    />
                                </span>
                            </div>
                        </div>
                    </Modal.Body>
                    <Modal.Footer style={style.footer} >
                        <button
                            type='button'
                            className='btn btn-primary btn-block btn-lg'
                            onClick={this.startOAuthConnectFlow}
                        >
                            {'Connect Webex'}
                        </button>

                        <button
                            type='button'
                            className='btn btn-link btn-sm'
                            onClick={this.handleClose}
                        >
                            {'Cancel'}
                        </button>

                    </Modal.Footer>
                </Modal>

            </div>
        );
    }
}

/* Define CSS styles here */
var getStyle = makeStyleFromTheme((theme) => {
    var xPos = ((window.innerWidth - 200) + 'px'); //shouldn't be set here as it doesn't rerender
    return {
        popover: {
            marginLeft: xPos,
            marginTop: '50px',
            maxWidth: '300px',
            height: '105px',
            width: '300px',
            background: theme.centerChannelBg,
        },
        popoverDM: {
            marginLeft: xPos,
            marginTop: '50px',
            maxWidth: '220px',
            height: '105px',
            width: '220px',
            background: theme.centerChannelBg,
        },
        header: {
            background: '#FFFFFF',
            color: '#0059A5',
            borderStyle: 'none',
            height: '10px',
            minHeight: '28px',
        },
        footer: {

            // margin: '0 auto',
            alignItems: 'center',
        },
        body: {
            padding: '0px 0px 10px 0px',
        },
        bodyText: {
            textAlign: 'center',
            margin: '20px 0 0 0',
            fontSize: '17px',
            lineHeight: '19px',
        },
        meetingId: {
            marginTop: '55px',
        },
        backdrop: {
            position: 'absolute',
            display: 'flex',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0, 0, 0, 0.50)',
            zIndex: 2000,
            alignItems: 'center',
            justifyContent: 'center',
        },
        modal: {

            // height: '350px',
            width: '400px',

            // position: 'relative',
            // margin: ''
            // margin: '30px auto',
            marginLeft: 'auto',
            marginRight: 'auto',

            // marginBottom: 'auto',
            maxWidth: '95%',

            // width: 600px;
            // margin-top: 30px;

        // display: 'flex !important',
        // alignItems: 'center',
        // position: 'absolute',
        // top: '50%',
        // left: '50%',
        // transform: 'translate(-50%, -50%) !important',
        // padding: '1em',
        // color: theme.centerChannelColor,
        // backgroundColor: theme.centerChannelBg,
        },
        modalCont: {
            maxWidth: '400px',
        },
        iconStyle: {
            position: 'relative',
            top: '-1px',
        },

        popoverBody: {
            maxHeight: '305px',
            overflow: 'auto',
            position: 'relative',
            width: '298px',
            left: '-14px',
            top: '-9px',
            borderBottom: '1px solid #D8D8D9',
        },

        popoverBodyDM: {
            maxHeight: '305px',
            overflow: 'auto',
            position: 'relative',
            width: '218px',
            left: '-14px',
            top: '-9px',
            borderBottom: '1px solid #D8D8D9',
        },
        logosConnect: {

        },
        mmLogo: {
            marginLeft: '40px',
        },
        webexLogo: {
            marginRight: '40px',
        },
        connectArrow: {
            marginTop: '34px',
        },
    };
});
