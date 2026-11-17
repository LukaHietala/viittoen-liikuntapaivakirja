import { createSignal, createResource } from "solid-js";
import ChallengesList from "./ChallengesList";

function CollectionItem(props) {
		const collection = props.collection;

		return (
				<>
						<div class="hierarchical-list">
								<div class="hierarchical-item">
										<div class="item-info">
												<strong>
														{collection.title}
												</strong>
										</div>
										<progress
												max={collection.goalPoints}
												value={collection.pot}
										>Edistys mittari</progress>
								</div>
								<ChallengesList
										challenges={() => collection.challenges}
										addPerformance={props.addPerformance}
								/>
						</div>
				</>
		);
}

export default CollectionItem;