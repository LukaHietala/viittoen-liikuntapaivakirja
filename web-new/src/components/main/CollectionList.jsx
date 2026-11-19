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
						<Show when={props.collections().length == 0}>
                <div class="block">
										<span>Ei päähaasteita</span>
                </div>
						</Show>
        </div>
    );
}

export default CollectionList;