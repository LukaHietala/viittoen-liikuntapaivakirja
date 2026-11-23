import CollectionItem from "./CollectionItem.jsx";

function CollectionList(props) {
    return (
        <div class="list">
						<For each={props.collections()}>{(item) =>
								<CollectionItem
										collection={item}
										editCollection={props.editCollection}
										deleteCollection={props.deleteCollection}
										
										addChallenge={props.addChallenge}
										editChallenge={props.editChallenge}
										deleteChallenge={props.deleteChallenge}
										
										deletePerformance={props.deletePerformance}
										users={props.users}
								/>
						}</For>
        </div>
    );
}

export default CollectionList;