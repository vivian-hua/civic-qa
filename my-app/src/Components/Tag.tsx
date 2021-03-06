import './Tag.css';

export type TagProps = {
    value: string;
    remove: Function;
}

export function Tag(props: TagProps) {
    return(
        <div className="tag">
            <p className="tag-value">{"#" + props.value}</p>
            <button className="tag-remove-btn" onClick={() => props.remove(props.value)}><img className="icon" src="./assets/icons/remove.svg"></img></button>
        </div>
    );
}