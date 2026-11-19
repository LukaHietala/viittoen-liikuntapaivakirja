import { createSignal } from "solid-js";
import { useAuth } from "../../context/AuthContext";
import ConfirmButton from "../common/ConfirmButton";

function PerformanceItem(props) {
		const { user } = useAuth();
		const performance = props.performance;
		
		const username = () => {
				const u = props.users().find((u)=> u.id === performance.userId);
				return u.name;
		};

		return (
				<div class="hierarchical-inline-item">
            <div>
								<span>
										Pisteet: <strong>{performance.points}</strong> -
								</span>
								<span>
										&nbsp;Suoritettu: <strong>{performance.createdAt}</strong> -
								</span>
								<span>
										&nbsp;Suorittanut: <strong>{username()}</strong>
								</span>
            </div>
						<div style="display: inline-flex; gap: 6px;">
								<Show
										when={user() && (user().isAdmin || user().id == performance.userId)}>
										<ConfirmButton
												onClick={() => props.deletePerformance(performance.id)}
										/>
								</Show>
						</div>
				</div>
		);
}

export default PerformanceItem;