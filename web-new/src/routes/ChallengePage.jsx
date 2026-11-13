import CollectionList from "../components/challenges/CollectionList";
import { createSignal, createResource } from "solid-js";
import CollectionForm from "../schemas/CollectionForm";
import { challengeService, collectionService, performanceService } from "../services";
import { useAuth } from "../context/AuthContext.jsx";

function ChallengePage() {
		const { user } = useAuth();
		const [adding, setAdding] = createSignal(false);
		
		const [collections, { refetch: refetchCollections }] = createResource(
				() => user() && user().id,
				async () => {
						return await collectionService.listAll();
				},
		);

		const addCollection = async (data) => {
				try {
						await collectionService.create(data);
				} catch (e) {
						console.error(e.message);
				}
				setAdding(false);
				// toast(`Lisättiin haaste: ${data.name}`)
				await refetchCollections();
		};

		const editCollection = async (data, id) => {
				try {
						await collectionService.updateByID(id, data);
				} catch (e) {
						console.error(e.message);
				}
				// toast(`Muokattiin haastetta: ${data.name}`)
				await refetchCollections();
		};

		const deleteCollection = async (id) => {
				try {
						await collectionService.deleteByID(id);
				} catch (e) {
						console.error(e.message);
				}
				// toast(`Poistettiin haaste: ${data.name}`)
				await refetchCollections();
		};

		const addChallenge = async (data) => {
				try {
						await challengeService.create(data);
				} catch (e) {
						console.error(e.message);
				}
				setAdding(false);
				// toast(`Lisättiin haaste: ${data.name}`)
				await refetchChallenges();
		};

		const editChallenge = async (data, id) => {
				try {
						await challengeService.updateByID(id, data);
				} catch (e) {
						console.error(e.message);
				}
				// toast(`Muokattiin haastetta: ${data.name}`)
				await refetchChallenges();
		};

		const deleteChallenge = async (id) => {
				try {
						await challengeService.deleteByID(id);
				} catch (e) {
						console.error(e.message);
				}
				// toast(`Poistettiin haaste: ${data.name}`)
				await refetchChallenges();
		};

		const deletePerformance = async (id) => {
				try {
						await performanceService.deleteByID(id);
				} catch (e) {
						console.error(e.message);
				}
				// toast(`Poistettiin suoritus: ${data.name}`)
				await refetchChallenges();
		};
		
		return (
				<>
						<h3>Haasteet</h3>
						<Show when={adding()}
									fallback={<button onClick={() => setAdding(true)}>Lisää isä haaste</button>}>
								<CollectionForm
										onSubmit={addCollection}
										collection={null}
										cancel={setAdding(false)}
								/>
						</Show>
						<Suspense fallback={<p>Ladataan...</p>}>
								<Show when={collections()}>
										<CollectionList
												collections={collections}
												editChallenge={editChallenge}
												deleteChallenge={deleteChallenge}
												deletePerformance={deletePerformance}
										/>
								</Show>
						</Suspense>
				</>
		);
}

export default ChallengePage;