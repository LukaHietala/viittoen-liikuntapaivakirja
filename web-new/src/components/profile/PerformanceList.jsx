import PerformanceItem from "./PerformanceItem.jsx";

function PerformanceList(props) {
    return (
        <div class="list">
            <For each={props.performances().toSorted(
										 (a, b) => a.createdAt < b.createdAt)}>{(item) =>
                <PerformanceItem
										performance={item}
										challenges={props.challenges}
										deletePerformance={props.deletePerformance}
								/>
										 }</For>
						<Show when={props.performances().length == 0}>
                <div class="item">
										Ei suorituksia
                </div>
						</Show>
        </div>
    );
}

export default PerformanceList;