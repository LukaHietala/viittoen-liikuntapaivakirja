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
        </div>
    );
}

export default PerformanceList;