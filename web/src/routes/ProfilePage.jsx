import PerformanceList from "../components/profile/PerformanceList";
import { createResource } from "solid-js";
import { performanceService, challengeService } from "../services";
import { useAuth } from "../context/AuthContext.jsx";
import { toast } from "../scripts/toast.js";

function ProfilePage() {
    const { user } = useAuth();
		
		const [performances, { refetch: refetchPerformances }] = createResource(
				() => user() && user().id,
				async () => {
						return await performanceService.listOwn();
				},
		);

		const [challenges] = createResource(
				async () => {
						return await challengeService.listAll();
				},
		);

		const deletePerformance = async (id) => {
				let success = true;
				try {
						await performanceService.deleteByID(id);
				} catch (e) {
						success = false;
						toast(e.message);
				}
				if (!success) return;
				toast(`Poistettiin suoritus`);			
				await refetchPerformances();
		};
		
		return (
				<>
						<Show when={user()}>
								<h2>Profiili ({user().name})</h2>
						</Show>
            <div class="block">
								<Show when={performances()}>
										<span>
												Keräämäsi pisteet: {performances().reduce(
														(sum, x) => sum + x.points, 0)}
										</span>
								</Show>
            </div>
						<h3>Suorituksesi</h3>
						<Show when={performances()}>
								<PerformanceList
										performances={performances}
										challenges={challenges}
										deletePerformance={deletePerformance}
								/>
						</Show>
				</>
		);
}

export default ProfilePage;