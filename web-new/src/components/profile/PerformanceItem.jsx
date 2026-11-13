import { createSignal } from "solid-js"

function PerformanceItem(props) {
		const performance = props.performance;
    const [showChallenge, setShowChallenge] = createSignal(false)

		return (
        <>
						{JSON.stringify(performance)}
            <Show when={showChallenge()}>
								haaste: potki tölkkejä
            </Show>
						<button onClick={() => setShowChallenge(!showChallenge())}>Näytä haaste</button>
						<button onClick={() => props.deletePerformance(performance.id)}>Poista</button>
        </>
		);
}

export default PerformanceItem;