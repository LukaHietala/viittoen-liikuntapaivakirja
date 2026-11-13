import ChallengesList from "../components/main/ChallengesList";
import { createResource } from "solid-js";
import { challengeService } from "../services";
import { performanceService } from "../services";
import { useAuth } from "../context/AuthContext.jsx";

function MainPage() {
		const { user } = useAuth();
		
		const [challenges, { refetch: refetchChallenges }] = createResource(
				() => user() && user().id,
				async () => {
						return await challengeService.listActive();
				},
		);

		const addPerformance = async (data) => {
				try {
						await performanceService.create(data);
				} catch (e) {
						console.error(e.message);
				}
				// toast(`Lisättiin haaste: ${data.name}`)
				await refetchChallenges();
		};
		
		return (
				<>
						<h3>Aktiiviset haasteet</h3>
						<Suspense fallback={<p>Ladataan...</p>}>
								<Show when={challenges()}>
										<ChallengesList
												challenges={challenges}
												addPerformance={addPerformance}
										/>
								</Show>
						</Suspense>
				</>
		);
}

export default MainPage;