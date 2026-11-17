import PerformanceItem from "./PerformanceItem.jsx";

function PerformanceList(props) {
    return (
				<div >
						<For each={props.performances().toSorted(
										 (a, b) => a.createdAt < b.createdAt)}>{(item) =>
												 <PerformanceItem
														 performance={item}
														 deletePerformance={props.deletePerformance}
												 />
										 }</For>
						<Show when={props.performances().length == 0}>
                <div class="hierarchical-item">
										Ei suorituksia
                </div>
						</Show>
				</div>
    );
}

export default PerformanceList;