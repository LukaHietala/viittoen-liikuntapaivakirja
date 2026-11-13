import { createSignal } from "solid-js"
import PerformanceForm from "../../schemas/PerformanceForm";

function ChallengesItem(props) {
		const challenge = props.challenge;
		
    const [show, setShow] = createSignal(false)

		const addT = (data) => {
				data.challengeID = challenge.id;
				props.addPerformance(data);
		};

		return (
        <>
						{challenge.title} - 
						{challenge.performances.reduce((sum, x) => sum + x.points, 0)}
            <Show when={show()}
									fallback={<button onClick={() => setShow(true)}>Lisää suoritus</button>}>
								<PerformanceForm
										onSubmit={addT}
										cancel={setShow(false)}
								/>
            </Show>
        </>
		);
}

export default ChallengesItem;