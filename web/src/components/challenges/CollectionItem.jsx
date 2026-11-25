import { createSignal, createResource } from "solid-js"
import PerformanceList from "./PerformanceList";
import ChallengeList from "./ChallengeList";
import CollectionForm from "../../schemas/CollectionForm";
import ChallengeForm from "../../schemas/ChallengeForm";
import ConfirmButton from "../common/ConfirmButton";
import { performanceService } from "../../services";
import { useAuth } from "../../context/AuthContext";
import penIcon from "../../assets/pen.svg"

function CollectionItem(props) {
		const { user } = useAuth();
		const collection = props.collection;

		// Helpers
		const pad = (x) => ("0" + x).slice(-2);
		const now = () => {
				const d = new Date();
				return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}`
		};
		const isActive = () => now() >= collection.startDate && now() <= collection.endDate;

		const [adding, setAdding] = createSignal(false);
		const [editing, setEditing] = createSignal(false);
		const [showing, setShowing] = createSignal(false);
		const [isConfirmingDelete, setIsConfirmingDelete] = createSignal(false);

		const editT = (data) => {
				props.editCollection(data, collection.id);
		};

		const addT = (data) => {
				props.addChallenge(data, collection.id);
		};

		return (
				<>
						<div class="hierarchical-list active">
								<div class="hierarchical-item">
										<div class="item-info">
                        <div>
														<span>{collection.title} {isActive() ? "(Aktiivinen)" : "(Ei aktiivinen)"}</span><br />
                            <span class="small">
																Pisteet: {collection.pot} / {collection.goalPoints} -
																Ajanjakso: {collection.startDate} -
																{collection.endDate}
														</span> 
                        </div>
										</div>
												<div class="item-actions">
														<Show when={!isConfirmingDelete()}>
																<button onClick={() => setShowing(!showing())}>
																		Näytä alihaasteet
																</button>
																<Show when={user() && user().isAdmin}>
																		<Show when={!editing()}>
																				<button class="btn-icon" onClick={() => setEditing(true)}>
																						<img src={penIcon} alt="Muokkaa" />
																				</button>	
																		</Show>
																</Show>
														</Show>
														<Show when={user() && (user().isAdmin || user().id == collection.userId)}>
																<ConfirmButton
																		onConfirmChange={setIsConfirmingDelete} 
																		onClick={() => props.deleteCollection(collection.id)}
																/>
														</Show>
												</div>
								</div>
								<Show when={editing()}>
										<CollectionForm
												title={`Muokataan: ${collection.title}`}
												collection={() => collection}
												onSubmit={editT}
												cancel={setEditing(false)}
										/>
								</Show>
								<Show when={showing()}>
										<Show when={adding()} fallback={
															<Show when={user() && user().isAdmin}>
																	<div class="hierarchical-item">
																			<div class="item-actions">
																					<button onClick={() => setAdding(!adding())}>
																							Lisää alihaaste
																					</button>
																			</div>
																	</div>
															</Show>
													}>
												<ChallengeForm
														title="Lisätään alihaaste"
                                                        description="Kun lisäät uuden alihaasteen siitä lähtee ilmotus muiden sähköpostiin."
														onSubmit={addT}
														challenge={null}
														cancel={setAdding(false)}
												/>
										</Show>
										
										<ChallengeList
												challenges={() => collection.challenges}
												editChallenge={props.editChallenge}
												deleteChallenge={props.deleteChallenge}
												deletePerformance={props.deletePerformance}
												users={props.users}
										/>
								</Show>
						</div>
				</>
		);
}

export default CollectionItem;