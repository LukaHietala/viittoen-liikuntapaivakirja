import ChallengesItem from "./ChallengesItem.jsx";

function ChallengesList(props) {
    return (
        <div>
            <For each={props.challenges()}>{(item) =>
                <ChallengesItem
										addPerformance={props.addPerformance}
										challenge={item}
										/>
            }</For>
						<Show when={props.challenges().length == 0}>
								Ei aktiivisia haasteita
						</Show>
        </div>
    );
}

export default ChallengesList;