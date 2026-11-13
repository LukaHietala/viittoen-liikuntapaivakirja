import PerformanceItem from "./PerformanceItem.jsx";

function PerformanceList(props) {
    return (
        <div>
            <For each={props.performances()}>{(item) =>
                <PerformanceItem
										performance={item}
										deletePerformance={props.deletePerformance}
								/>
            }</For>
						<Show when={props.performances().length == 0}>
								Ei suorituksia
						</Show>
        </div>
    );
}

export default PerformanceList;