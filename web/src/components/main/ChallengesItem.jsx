import { createSignal } from "solid-js"
import PerformanceForm from "../../schemas/PerformanceForm";

function ChallengesItem(props) {
		const challenge = props.challenge;
		const ownerColor = () => {
				const u = props.users().find((u)=> u.id === challenge.userId);
				return u ? u.color : "black";
		};
		// Helpers
		const pad = (x) => ("0" + x).slice(-2);
		const now = () => {
				const d = new Date();
				return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}`
		};
		const isActive = () => now() > challenge.startDate && now() < challenge.endDate;
    const [show, setShow] = createSignal(false);

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
						<div class="hierarchical-item"
								 class={isActive() ? "" : "disactive"}>
								<div class="item-info">
										<div class="laser" style={{"background-color": ownerColor()}}></div>
										<span>
												{challenge.title} - Pisteet: {
														challenge.performances.reduce((sum, x) => sum + x.points, 0)}
										</span>
								</div>
								<div class="item-actions">
										<button onClick={() => setShow(true)}
														disabled={!isActive()}>
												Lisää suoritus
										</button>
								</div>
						</div>
        </Show>
		);
}

export default ChallengesItem;