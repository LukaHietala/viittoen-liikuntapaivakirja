import ChallengesItem from "./ChallengesItem.jsx";

function ChallengesList(props) {
    return (
				<>
						<For each={props.challenges()}>{(item) =>
								<ChallengesItem
										addPerformance={props.addPerformance}
										challenge={item}
										users={props.users}
								/>
						}</For>
						<Show when={props.challenges().length == 0}>
                <div class="hierarchical-item">
										Ei aktiivisia haasteita
                </div>
						</Show>
				</>
    );
}

export default ChallengesList;