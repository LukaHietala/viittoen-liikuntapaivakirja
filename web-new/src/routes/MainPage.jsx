import CollectionList from "../components/main/CollectionList";
import { createResource } from "solid-js";
import { collectionService, performanceService } from "../services";
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

		const addPerformance = async (data) => {
				try {
						await performanceService.create(data);
				} catch (e) {
						console.error(e.message);
				}
				toast("Lisättiin suoritus");
				await refetchCollections();
		};
		
		return (
				<>
						<h2>Aktiiviset haasteet</h2>
						<Suspense fallback={<p>Ladataan...</p>}>
								<Show when={collections()}>
										<CollectionList
												collections={collections}
												addPerformance={addPerformance}
										/>
								</Show>
						</Suspense>
				</>
		);
}

export default MainPage;