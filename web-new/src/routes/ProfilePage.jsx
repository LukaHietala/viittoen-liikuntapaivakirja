import PerformanceList from "../components/profile/PerformanceList";
import { createResource } from "solid-js";
//import { challengeService } from "../services";

function ProfilePage() {
		// const [performances] = createResource(
		// 		// () => user() && user().id,
		// 		async () => {
		// 				return await challengeService.listActive();
		// 		},
		// );

		const deletePerformance = async (id) => {
				try {
						console.log(id);
						//await ??
				} catch (e) {
						console.error(e.message);
				}
				// toast(`Poistettiin suoritus: ${data.name}`)
				await refetch();
		};
		
		return (
				<>
						<h3>Profiili (pedofiili)</h3>
						<p>Keräämäsi pisteet:</p>
						<p>Suorituksesi:</p>
						<PerformanceList
								performances={() => [{id:1},{id:2},{id:3}]}
								deletePerformance={deletePerformance}
						/>
				</>
		);
}

export default ProfilePage;