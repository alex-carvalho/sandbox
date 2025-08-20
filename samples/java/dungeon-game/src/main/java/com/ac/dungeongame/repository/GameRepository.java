package com.ac.dungeongame.repository;

import com.ac.dungeongame.model.Game;
import org.springframework.data.jpa.repository.JpaRepository;

public interface GameRepository extends JpaRepository<Game, Long> {
}
