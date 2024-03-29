/*
 * This Java source file was generated by the Gradle 'init' task.
 */
package com.ac;

import org.apache.pekko.actor.*;

import java.util.concurrent.ExecutionException;
import java.util.stream.IntStream;

public class App {


    public static void main(String[] args) throws ExecutionException, InterruptedException {
        ActorSystem system = ActorSystem.create("test-system");

        ActorRef actorRef = system.actorOf(Props.create(SenderActor.class), "string-sender");
        IntStream.range(0, 10).parallel().forEach(value ->  {
            actorRef.tell(String.valueOf(value), ActorRef.noSender());
        });


    }



    public static class SenderActor extends AbstractLoggingActor {

        final ActorRef actorRef;
        public SenderActor() {
            actorRef = getContext().actorOf(Props.create(WordLengthActor.class), "word-counter-actor");
        }

        @Override
        public Receive createReceive() {
            return receiveBuilder()
                    .match(String.class, i -> {
                        try {
                            log().info("Received message from " + getSender() + " with value: " + i);
                            actorRef.tell(new SizeText("this is a text"), getSelf());
                        } catch (Exception ex) {
                            getSender().tell(new FSM.Failure(ex), getSelf());
                            throw ex;
                        }
                    })
                    .match(Integer.class, r -> {
                        log().info("Received message from " + getSender() + " with value: " + r);
                    })
                    .build();
        }
    }
    public static class WordLengthActor extends AbstractLoggingActor {


        @Override
        public Receive createReceive() {
            return receiveBuilder()
                    .match(SizeText.class, r -> {
                        try {
                            log().info("Received message from " + getSender());
                            int size = r.line.length();
                            getSender().tell(size, getSelf());
                        } catch (Exception ex) {
                            getSender().tell(new FSM.Failure(ex), getSelf());
                            throw ex;
                        }
                    })
                    .build();
        }

    }

    public static final class SizeText {
        String line;

        public SizeText(String line) {
            this.line = line;
        }
    }

}
