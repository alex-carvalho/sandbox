package com.ac.dungeongame.model;

import jakarta.persistence.*;

@Entity
public class AlgorithmExecution {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id;

    @ManyToOne
    @JoinColumn(name = "game_id", nullable = false)
    private Game game;

    private String algorithmType;

    private long durationInNanos;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public Game getGame() {
        return game;
    }

    public void setGame(Game game) {
        this.game = game;
    }

    public String getAlgorithmType() {
        return algorithmType;
    }

    public void setAlgorithmType(String algorithmType) {
        this.algorithmType = algorithmType;
    }

    public long getDurationInNanos() {
        return durationInNanos;
    }

    public void setDurationInNanos(long durationInNanos) {
        this.durationInNanos = durationInNanos;
    }
}
