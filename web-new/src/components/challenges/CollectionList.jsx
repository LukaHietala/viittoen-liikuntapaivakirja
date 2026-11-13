import CollectionItem from "./CollectionItem.jsx";

function CollectionList(props) {
    return (
        <div>
						<For each={props.collections()}>{(item) =>
								<CollectionItem
										collection={item}
										editCollection={props.editCollection}
										deleteCollection={props.deleteCollection}
										deletePerformance={props.deletePerformance}
								/>
						}</For>
						<Show when={props.collections().length == 0}>
								Ei haasteita
						</Show>
        </div>
    );
}

export default CollectionList;