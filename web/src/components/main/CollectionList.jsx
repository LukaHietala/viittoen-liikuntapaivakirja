import CollectionItem from "./CollectionItem.jsx";

function CollectionList(props) {
    return (
        <div class="list">
						<For each={props.collections()}>{(item) =>
								<CollectionItem
										collection={item}
										addPerformance={props.addPerformance}
										users={props.users}
								/>
						}</For>	
        </div>
    );
}

export default CollectionList;