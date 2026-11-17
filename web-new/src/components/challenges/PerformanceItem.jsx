import { createSignal } from "solid-js";
import { useAuth } from "../../context/AuthContext";
import ConfirmButton from "../common/ConfirmButton";

function PerformanceItem(props) {
		const { user } = useAuth();
		const performance = props.performance;

		return (
						<div class="hierarchical-item">
            <div class="item-info">
								{performance.points} - {performance.createdAt}
            </div>
        <div class="item-actions">
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