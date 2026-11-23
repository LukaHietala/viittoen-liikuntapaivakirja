import ChallengesList from "./ChallengesList";

function CollectionItem(props) {
		const collection = props.collection;

		return (
				<>
						<div class="hierarchical-list active">
								<div class="hierarchical-item">
										<div class="item-info">
												<div>
														<strong>
																{collection.title}
														</strong>
														<br/>
														<span class="small">
																Pisteet: {collection.pot} / {collection.goalPoints} -
																Ajanjakso: {collection.startDate} - {collection.endDate}
														</span>
												</div>
										</div>
										<progress
												max={collection.goalPoints}
												value={collection.pot}
										>Edistys mittari</progress>
								</div>
								<ChallengesList
										challenges={() => collection.challenges}
										addPerformance={props.addPerformance}
										users={props.users}
								/>
						</div>
				</>
		);
}

export default CollectionItem;