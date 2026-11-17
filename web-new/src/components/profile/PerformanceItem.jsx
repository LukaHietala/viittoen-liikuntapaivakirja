import { createSignal } from "solid-js";
import { performanceService } from "../../services";
import trashCan from "../../assets/can.svg"
import ConfirmButton from "../common/ConfirmButton";

function PerformanceItem(props) {
		const performance = props.performance;
    const [showChallenge, setShowChallenge] = createSignal(false)

    const challenge =
					() => props.challenges().find((x) => x.id == performance.challengeId);

		return (
        <>
            <div class="item">
                <div class="item-info">
                    <span>
                        Pisteet: <strong>{performance.points}</strong>
                    </span>
                    <span>
                        - Suoritettu <strong>{performance.createdAt}</strong>
                    </span>
                </div>
                <div class="item-actions">
                    <button onClick={() => setShowChallenge(!showChallenge())}>
                        Näytä haaste
                    </button>
										<ConfirmButton
												onClick={() => props.deletePerformance(performance.id)}
										/>
                </div>
            </div>
            <Show when={showChallenge() && challenge()}>
                <div class="item">
										<div>
                        <strong>{challenge().title}</strong><br/>
                        <span>
														Ajanjakso: {challenge().startDate} - {challenge().endDate}
												</span>
										</div>
                </div>
            </Show>
        </>
		);
}

export default PerformanceItem;