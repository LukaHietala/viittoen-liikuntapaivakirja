import CollectionList from "../components/main/CollectionList";
import { createResource } from "solid-js";
import { collectionService, performanceService, userService } from "../services";
import { useAuth } from "../context/AuthContext.jsx";
import { toast } from "../scripts/toast.js";

function MainPage() {
		const { user } = useAuth();
		
		const [collections, { refetch: refetchCollections }] = createResource(
				() => user() && user().id,
				async () => {
						return await collectionService.listActive();
				},
		);

		const [users] = createResource(
				() => user() && user().id,
				async () => {
						return await userService.listAll();
				},
		);


		const addPerformance = async (data) => {
				try {
						await performanceService.create(data);
				} catch (e) {
						toast(e.message);
				}
				toast("Lisättiin suoritus");
				await refetchCollections();
		};
		
		return (
				<>
						<h2>Aktiiviset haasteet</h2>
						<Suspense fallback={<p>Ladataan...</p>}>
								<Show when={collections() && users()}>
										<CollectionList
												collections={collections}
												addPerformance={addPerformance}
												users={users}
										/>
								</Show>
						</Suspense>
				</>
		);
}

export default MainPage;