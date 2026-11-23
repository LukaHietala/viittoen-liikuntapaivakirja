import PerformanceItem from "./PerformanceItem.jsx";

function PerformanceList(props) {
    return (
				<div >
						<For each={props.performances().toSorted(
										 (a, b) => a.createdAt < b.createdAt)}>{(item) =>
												 <PerformanceItem
														 performance={item}
														 deletePerformance={props.deletePerformance}
														 users={props.users}
												 />
										 }</For>
						<Show when={props.performances().length == 0}>
                <div style="padding: 15px;">
										Ei suorituksia
                </div>
						</Show>
				</div>
    );
}

export default PerformanceList;