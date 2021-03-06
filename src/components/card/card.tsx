import React from 'react'
import { LiveStreams } from '../main/main'
import { Button } from '../styled_comp/styles'
import { Visit, VisitSpan } from './card_styles'

interface Props {
    streamerData: LiveStreams;
    setSelected: Function;
    setDetails: Function;

}
const Card = (props: Props) => {
    const { thumbnails, likes, dislikes, viewers, title, channelId, description, imageId, videoId, type, name, displayName, isPlaying } = props.streamerData
    const newTitle = title.slice(0, 24)
    const image: string = imageId.startsWith("https") ? imageId : `https://s3.us-east-2.amazonaws.com/xhnetwork/${imageId}`
    const theme = {
        marginLeft: "15px",
        isYoutube: type === "youtube" ? true : false
    }
    return (
        <div className="card">
            <Visit
                theme={theme}
                target="_blank"
                rel="noopener noreferrer"
                href={type === "youtube" ? `https://www.youtube.com/watch?v=${channelId}` : `https://www.twitch.tv/${name}`}
            >
                <VisitSpan>
                    Visit Channel
                </VisitSpan>
            </Visit>
            {likes && (
                <div className="likes">
                    <span>
                        <i className="fas fa-thumbs-up" /> {likes}
                    </span>
                    <span>
                        <i className="fas fa-thumbs-down" /> {dislikes}
                    </span>
                </div>
            )}
            <div className="image">
                <img src={thumbnails.high || thumbnails.low} alt="thumbnail" />
            </div>
            <div className="details">
                <div className="content">
                    <img src={image} alt="streamer" style={type !== "youtube" ? { border: "3px solid #4B367C" } : { border: "solid 3px red " }} />
                    <div className="content-details">
                        <h3>{newTitle}</h3>
                        <span style={{ margin: '-15px auto 0 auto' }}>{displayName || name} on {type}</span>
                        {isPlaying && (
                            <span><small>{isPlaying}</small></span>
                        )}
                        <span><i style={{ color: "red" }} className="fas fa-dot-circle" /> {viewers + " viewers"}</span>
                    </div>
                </div>
                <div className="buttons">
                    <Button
                        theme={theme}
                        onClick={() => {
                            if (window.innerWidth <= 900) {
                                const youtubeLink: string = type === "youtube" ? `https://www.youtube.com/watch?v=${videoId}` : `https://www.twitch.tv/${name}`;
                                const win: any = window.open(youtubeLink, '_blank');
                                win.focus();
                                return;
                            }
                            props.setSelected(channelId)
                        }}
                    >Watch</Button>
                        <span className="details"
                            onClick={() => {
                                props.setDetails(props.streamerData)
                            }}
                        >
                            More Details
                    </span>
                </div>
            </div>
        </div>
    )
}

export default Card