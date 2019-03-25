// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import React from 'react';
import PropTypes from 'prop-types';

import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

import {Svgs} from '../../constants';
import {formatDate} from '../../utils/date_utils';

export default class PostTypeWebex extends React.PureComponent {
    static propTypes = {

        /*
         * The post to render the message for.
         */
        post: PropTypes.object.isRequired,

        /**
         * Set to render post body compactly.
         */
        compactDisplay: PropTypes.bool,

        /**
         * Flags if the post_message_view is for the RHS (Reply).
         */
        isRHS: PropTypes.bool,

        /**
         * Set to display times using 24 hours.
         */
        useMilitaryTime: PropTypes.bool,

        /*
         * Logged in user's theme.
         */
        theme: PropTypes.object.isRequired,

        /*
         * Creator's name.
         */
        creatorName: PropTypes.string.isRequired,
    };

    static defaultProps = {
        mentionKeys: [],
        compactDisplay: false,
        isRHS: false,
    };

    render() {
        const style = getStyle(this.props.theme);
        const post = this.props.post;
        const props = post.props || {};

        let preText;
        let content;
        let subtitle;
        if (props.meeting_status === 'STARTED') {
            preText = `${this.props.creatorName} has started a meeting`;
            content = (
                <a
                    className='btn btn-lg btn-primary'
                    style={style.button}
                    rel='noopener noreferrer'
                    target='_blank'
                    href={props.meeting_link}
                >
                    <i
                        style={style.buttonIcon}
                    >
                        <svg
                            width='24px'
                            height='24px'
                            viewBox='0 0 24 24'
                            style={style.webexLogo}
                        >
                            <path
                                fill='#FFFFFF'
                                d='M12,3A9,9 0 0,1 21,12A9,9 0 0,1 12,21A9,9 0 0,1 3,12A9,9 0 0,1 12,3M5.94,8.5C4,11.85 5.15,16.13 8.5,18.06C11.85,20 18.85,7.87 15.5,5.94C12.15,4 7.87,5.15 5.94,8.5Z'
                            />
                        </svg>
                    </i>
                    {'JOIN MEETING'}
                </a>
            );

            if (props.meeting_personal) {
                subtitle = (
                    <span>
                        {'Personal Meeting ID (PMI) : '}
                        <a
                            rel='noopener noreferrer'
                            target='_blank'
                            href={props.meeting_link}
                        >
                            {props.meeting_id}
                        </a>
                    </span>
                );
            } else {
                subtitle = (
                    <span>
                        {'Meeting ID : '}
                        <a
                            rel='noopener noreferrer'
                            target='_blank'
                            href={props.meeting_link}
                        >
                            {props.meeting_id}
                        </a>
                    </span>
                );
            }
        }

        let title = 'Webex Meeting';
        if (props.meeting_topic) {
            title = props.meeting_topic;
        }

        return (
            <div>
                {preText}
                <div style={style.attachment}>
                    <div style={style.content}>
                        <div style={style.container}>
                            <h1 style={style.title}>
                                {title}
                            </h1>
                            {subtitle}
                            <div>
                                <div style={style.body}>
                                    {content}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

const getStyle = makeStyleFromTheme((theme) => {
    return {
        attachment: {
            marginLeft: '-20px',
            position: 'relative',
        },
        content: {
            borderRadius: '4px',
            borderStyle: 'solid',
            borderWidth: '1px',
            borderColor: '#BDBDBF',
            margin: '5px 0 5px 20px',
            padding: '2px 5px',
        },
        container: {
            borderLeftStyle: 'solid',
            borderLeftWidth: '4px',
            padding: '10px',
            borderLeftColor: '#89AECB',
        },
        body: {
            overflowX: 'auto',
            overflowY: 'hidden',
            paddingRight: '5px',
            width: '100%',
        },
        title: {
            fontSize: '16px',
            fontWeight: '600',
            height: '22px',
            lineHeight: '18px',
            margin: '5px 0 1px 0',
            padding: '0',
        },
        button: {
            fontFamily: 'Open Sans',
            fontSize: '12px',
            fontWeight: 'bold',
            letterSpacing: '1px',
            lineHeight: '19px',
            marginTop: '12px',
            borderRadius: '4px',

            // maxHeight: '20px',
            paddingTop: '0px',
            paddingLeft: '12px',
            color: theme.buttonColor,
        },
        buttonIcon: {
            padding: '0',

            marginTop: '5px',
            paddingRight: '6px',
            width: '14px',
            height: '19px',
            fill: theme.buttonColor,
        },
        summary: {
            fontFamily: 'Open Sans',
            fontSize: '14px',
            fontWeight: '600',
            lineHeight: '26px',
            margin: '0',
            padding: '14px 0 0 0',
        },
        summaryItem: {
            fontFamily: 'Open Sans',
            fontSize: '14px',
            lineHeight: '26px',
        },
        webexLogo: {
            width: '24px',
            height: '24px',
            position: 'relative',
            top: '7px',
        }
    };
});
