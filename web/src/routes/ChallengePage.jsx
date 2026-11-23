import CollectionList from "../components/challenges/CollectionList";
import { createSignal, createResource } from "solid-js";
import CollectionForm from "../schemas/CollectionForm";
import { challengeService, collectionService, performanceService, userService } from "../services";
import { useAuth } from "../context/AuthContext.jsx";
import { toast } from "../scripts/toast.js";

function ChallengePage() {
		const { user } = useAuth();
		const [adding, setAdding] = createSignal(false);
		
		const [collections, { refetch: refetchCollections }] = createResource(
				() => user() && user().id,
				async () => {
						return await collectionService.listAll();
				},
		);

		const [users] = createResource(
				() => user() && user().id,
				async () => {
						return await userService.listAll();
				},
		);

		const addCollection = async (data) => {
				let success = true;
				try {
						await collectionService.create(data);
				} catch (e) {
						success = false;
						toast(e.message);
				}
				if (!success) return;
				setAdding(false);
				toast(`Lisättiin haaste: ${data.title}`)				
				await refetchCollections();
		};

		const editCollection = async (data, id) => {
				let success = true;
				try {
						await collectionService.updateByID(id, data);
				} catch (e) {
						success = false;
						toast(e.message);
				}
				if (!success) return;
				toast(`Muokattiin haastetta: ${data.title}`)
				await refetchCollections();
		};

		const deleteCollection = async (id) => {
				let success = true;
				try {
						await collectionService.deleteByID(id);
				} catch (e) {
						success = false;
						toast(e.message);
				}
				if (!success) return;
				toast("Poistettiin haaste")
				
				await refetchCollections();
		};

		const addChallenge = async (data, id) => {
				let success = true;
				try {
						data.collectionId = id
						await challengeService.create(data);
				} catch (e) {
						success = false;
						toast(e.message);
				}
				if (!success) return;
				setAdding(false);
				toast(`Lisättiin haaste: ${data.title}`)
				
				await refetchCollections();
		};

		const editChallenge = async (data, id) => {
				let success = true;
				try {
						await challengeService.updateByID(id, data);
				} catch (e) {
						success = false;
						toast(e.message);
				}
				if (!success) return;
				toast(`Muokattiin haastetta: ${data.title}`)
				
				await refetchCollections();
		};

		const deleteChallenge = async (id) => {
				let success = true;
				try {
						await challengeService.deleteByID(id);
				} catch (e) {
						success = false;
						toast(e.message);
				}
				if (!success) return;
				toast(`Poistettiin haaste`)
				
				await refetchCollections();
		};

		const deletePerformance = async (id) => {
				let success = true;
				try {
						await performanceService.deleteByID(id);
				} catch (e) {
						success = false;
						toast(e.message);
				}
				if (!success) return;
				toast("Poistettiin suoritus")
				
				await refetchCollections();
		};
		
		return (
				<>
						<h2>Haasteet</h2>
            <Show when={adding()}
                  fallback={
                      <div class="block">
                          <button onClick={() => setAdding(true)}>
                              Luo päähaaste
                          </button>
                      </div>
                  }>
                <CollectionForm
                    onSubmit={addCollection}
                    collection={null}
                    cancel={setAdding(false)}
                />
            </Show>
            <h3>Päähaasteet</h3>
						<Suspense fallback={<p>Ladataan...</p>}>
								<Show when={collections() && users()}>
										<CollectionList
												collections={collections}
												editCollection={editCollection}
												deleteCollection={deleteCollection}
												
												addChallenge={addChallenge}
												editChallenge={editChallenge}
												deleteChallenge={deleteChallenge}
												
												deletePerformance={deletePerformance}
												users={users}
										/>
								</Show>

                <Show when={collections()?.length == 0}>
                    <div class="block">
                        Ei päähaasteita
                    </div>
                </Show>
						</Suspense>
				</>
		);
}

export default ChallengePage;