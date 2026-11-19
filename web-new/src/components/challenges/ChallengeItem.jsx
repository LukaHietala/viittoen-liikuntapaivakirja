import { createSignal, createResource } from "solid-js";
import PerformanceList from "./PerformanceList";
import ChallengeForm from "../../schemas/ChallengeForm";
import ConfirmButton from "../common/ConfirmButton";
import { performanceService } from "../../services";
import { useAuth } from "../../context/AuthContext";
import penIcon from "../../assets/pen.svg"

function ChallengeItem(props) {
		const { user } = useAuth();
		const challenge = props.challenge;
		const ownerColor = () => {
				const u = props.users().find((u)=> u.id === challenge.userId);
				return u ? u.color : "black";
		};
		
		const [editing, setEditing] = createSignal(false)
		const [showPerformances, setShowPerformances] = createSignal(false)

		const editT = (data) => {
				props.editChallenge(data, challenge.id);
		};

		return (
        <>
						<div class="hierarchical-item">
								<div class="item-info">										
										<div class="laser" style={{"background-color": ownerColor()}}></div>
                    <div>
                        {/*TODO, TILA: Aktiivinen/vanhentunut/odottaa*/}
                        <div>
                            <span>{challenge.title}</span><br />
                            <span class="small">
																Pisteet: { challenge.performances.
																					 reduce((sum, x) => sum + x.points, 0)
																				 }  /
																{challenge.goalPoints} -
																Ajanjakso: {challenge.startDate} - {challenge.endDate}
														</span>
                        </div>
                    </div>
								</div>
								<div class="item-actions">
										<button onClick={() => setShowPerformances(!showPerformances())}>
												Näytä suoritukset
										</button>
										<Show when={!editing()}>
												<button class="btn-icon" onClick={() => setEditing(true)}>
														<img src={penIcon} alt="Muokkaa" />
												</button>
										</Show>
                    {/* TODO: Piilota muut napit kun confirm tilassa */}
										<Show
												when={user() && (user().isAdmin || user().id == challenge.userId)}>
												<ConfirmButton
														onClick={() => props.deleteChallenge(challenge.id)}
												/>
										</Show> 
										
								</div>
						</div>
						<Show when={editing()}>
								<ChallengeForm
										title={`Muokataan: ${challenge.title}`}
										challenge={() => challenge}
										onSubmit={editT}
										cancel={setEditing(false)}
								/>
						</Show>
						<Show when={showPerformances()}>
								<PerformanceList
										performances={() => challenge.performances}
										deletePerformance={props.deletePerformance}
										users={props.users}
								/>
						</Show>
        </>
				
		);
}

export default ChallengeItem;