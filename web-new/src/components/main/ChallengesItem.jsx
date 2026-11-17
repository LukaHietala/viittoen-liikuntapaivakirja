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
            <Show when={!show()} fallback={
                <PerformanceForm
                title={challenge.title}
                description={challenge.description}
                unit={challenge.unit}
                onSubmit={addT}
                cancel={setShow(false)}
            />

            }>
				<div class="hierarchical-item">
						<div class="item-info">
								{challenge.title} -
								Edistys:
								<progress
										max={challenge.goalPoints}
										value={challenge.performances.reduce((sum, x) => sum + x.points, 0)}
								>

										Edistys mittari
								</progress>
						</div>
						<div class="item-actions">
            <button onClick={() => setShow(true)}>
            Lisää suoritus
            </button>
						</div>
        </div>
            </Show>
		);
}

export default ChallengesItem;