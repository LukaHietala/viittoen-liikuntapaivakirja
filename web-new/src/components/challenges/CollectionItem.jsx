import { createSignal, createResource } from "solid-js"
import PerformanceList from "./PerformanceList";
import ChallengeList from "./ChallengeList";
import ChallengeForm from "../../schemas/ChallengeForm";
import { performanceService } from "../../services";

function CollectionItem(props) {
		const collection = props.collection;

		const [adding, setAdding] = createSignal(false);
    const [editing, setEditing] = createSignal(false)
		const [deleting, setDeleting] = createSignal(false);

		const editT = (data) => {
				props.editCollection(data, challenge.id);
		};

		return (
        <>
						{collection.title}
            <Show when={editing()}
									fallback={<button onClick={() => setEditing(true)}>Muokkaa</button>}>
								moi
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
								<button onClick={() => props.deleteCollection(challenge.id)}>
										Varmista
								</button>
						</Show>
						Haasteet:
						<Show when={adding()}
									fallback={
											<button onClick={() => setAdding(true)}>
													Lisää haaste
											</button>
									}>
								<ChallengeForm
										onSubmit={props.addChallenge}
										challenge={null}
										cancel={setAdding(false)}
								/>
						</Show>
						<ChallengeList
								challenges={() => collection.challenges}
								editChallenge={props.editChallenge}
								deleteChallenge={props.deleteChallenge}
								deletePerformance={props.deletePerformance}
						/>
        </>
		);
}

export default CollectionItem;