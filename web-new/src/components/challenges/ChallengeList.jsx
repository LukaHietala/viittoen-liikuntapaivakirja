import ChallengeItem from "./ChallengeItem.jsx";

function ChallengeList(props) {
    return (
        <div>
						<For each={props.challenges()}>{(item) =>
								<ChallengeItem
										challenge={item}
										editChallenge={props.editChallenge}
										deleteChallenge={props.deleteChallenge}
										deletePerformance={props.deletePerformance}
								/>
						}</For>
						<Show when={props.challenges().length == 0}>
								Ei haasteita
						</Show>
        </div>
    );
}

export default ChallengeList;