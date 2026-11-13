import { createSignal, createResource } from "solid-js"
import PerformanceList from "./PerformanceList";
import ChallengeForm from "../../schemas/ChallengeForm";
import { performanceService } from "../../services";

function ChallengeItem(props) {
		const challenge = props.challenge;
		
    const [editing, setEditing] = createSignal(false)
		const [deleting, setDeleting] = createSignal(false);
		const [showPerformances, setShowPerformances] = createSignal(false)

		const editT = (data) => {
				props.editChallenge(data, challenge.id);
		};

		return (
        <>
						{challenge.title}
            <Show when={editing()}
									fallback={<button onClick={() => setEditing(true)}>Muokkaa</button>}>
								<ChallengeForm
										challenge={() => challenge}
										onSubmit={editT}
										cancel={setEditing(false)}
								/>
            </Show>
						<Show when={deleting()}
									fallback={
											<button onClick={() => setDeleting(true)}>
													Poista
											</button>
									}>
								<button onClick={() => setDeleting(false)}>
										Peru
								</button>
								<button onClick={() => props.deleteChallenge(challenge.id)}>
										Varmista
								</button>
						</Show>
						<button onClick={() => setShowPerformances(!showPerformances())}>Näytä suoritukset</button>
						<Show when={showPerformances()}>
								<PerformanceList
										performances={() => challenge.performances}
										deletePerformance={props.deletePerformance}
								/>
						</Show>
        </>
		);
}

export default ChallengeItem;