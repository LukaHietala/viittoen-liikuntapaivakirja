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
				let success = true;
				try {
						await performanceService.create(data);
				} catch (e) {
						success = false;
						toast(e.message);
				}
				if (!success) return;
				toast("Lisättiin suoritus onnistuneesti");			
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
                <Show when={collections()?.length == 0}>
                    <div class="block">
                        <span>Ei aktiivisia päähaasteita</span>
                    </div>
                </Show>
            </Suspense>
				</>
		);
}

export default MainPage;