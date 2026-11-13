import { createSignal } from "solid-js"

function PerformanceItem(props) {
		const performance = props.performance;

		return (
        <>
						{JSON.stringify(performance)}
						<button onClick={() => props.deletePerformance(performance.id)}>Poista</button>
        </>
		);
}

export default PerformanceItem;